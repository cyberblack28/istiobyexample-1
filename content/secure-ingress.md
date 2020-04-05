---
title: "セキュアなIngress"
publishDate: "2019-12-31"
categories: ["Security"]
---

Kubernetesクラスター上でワークロードを実行している場合、その一部をクラスターの外部に公開したくなることがあるでしょう。[Istio Ingress Gateway](/ingress)は、1つまたは複数のバックエンドホストの内向けトラフィックをルーティングできるカスタマイズ可能なプロキシです。しかし、HTTPSを用いたセキュアなIngressトラフィックを受け付けるにはどうすれば良いでしょうか？

Istioは、証明書と鍵をIngress GatewayにマウントすることでTLS Ingressをサポートし、内向けのトラフィックをクラスター内サービスに安全にルーティングできるようにします。Istioでセキュアなイングレスを設定すると、Ingress GatewayがすべてのTLS操作（ハンドシェイク、証明書/鍵交換）を処理し、アプリケーションコードからTLSを切り離すことができます。さらに、TLSトラフィックにIngress Gatewayを使用すると、組織全体の証明書と鍵の管理を一元化および自動化できます。

Istioは、2つの方法によるIngress Gatewayの保護をサポートしています。 1つは[ファイルマウント](https://istio.io/docs/tasks/traffic-management/ingress/secure-ingress-mount/)による方法で、IngressGatewayの証明書と鍵を生成し、KubernetesのSecretとして手動でIngressGatewayにマウントします。2つ目の方法は、IngressGateway PodでIstioプロキシと一緒に実行されるエージェントである[Secrets Discovery Service](https://istio.io/docs/tasks/traffic-management/ingress/secure-ingress-sds/)（SDS）を使用する方法です。SDSエージェントは istio-system ネームスペースを監視して新しいシークレットを探し、ユーザーに代わってそれらをゲートウェイのプロキシにマウントします。ファイルマウント方式と同様に、SDSはサーバーサイドと相互TLSの両方をサポートします。

SDSメソッドを使用して、相互HTTPS認証でIngress Gatewayを構成する方法を見てみましょう。

![](/images/secure-ingress-arch.png)

ここでは、FooCorpと呼ばれる建設資材企業が1つのKubernetesクラスターを運用しています。`ux` という1つのチームが、顧客向けのWeb `frontend` サービスを実行しています。もう1つは `corp-services` で、内部向けの `inventory` サービスを実行してサプライチェーンを追跡します。どちらのサービスにも独自の `foocorp` サブドメインがあり、セキュリティチームはすべてのサービスに独自の証明書と鍵を持たせることを義務付けています。このクラスターでのセキュアな入力の構成を見ていきましょう。

まず、[global SDS ingress](https://istio.io/docs/reference/config/installation-options/#gateways-options)オプションを有効にしてIstioをインストールします。（SDS ingressの有効化はインストールオプションです）これを有効にすると、Istio `ingress-gateway` Podには、2つのコンテナー、`istio-proxy`（Envoy）と `ingress-sds`（Secrets Discoveryエージェント）が立てられます。:

```
istio-ingressgateway-6f7d65d984-m2zmn     2/2     Running     0          44s
```

次に、`ux` と `corp-services` という2つのネームスペースを作成し、両方にIstioサイドカープロキシインジェクション用のラベルを付けます。次に、`frontend`（`ux` namespace）と`inventory`（`corp-services` namespace）にDeploymentとServiceを作成します。

これで、`frontend.foocorp.com` と `inventory.foocorp.com` という2つのドメインの証明書と鍵を生成する準備が整いました。（注：これを試すためにドメイン名を購入する必要はありません。以降のステップでは、`host` ヘッダーを使用してテストします。）これらの鍵からKubernetes Secretを生成します。:

```
kubectl create -n istio-system secret generic frontend-credential  \
--from-file=key=frontend.foocorp.com/3_application/private/frontend.foocorp.com.key.pem \
--from-file=cert=frontend.foocorp.com/3_application/certs/frontend.foocorp.com.cert.pem \
--from-file=cacert=frontend.foocorp.com/2_intermediate/certs/ca-chain.cert.pem

kubectl create -n istio-system secret generic inventory-credential  \
--from-file=key=inventory.foocorp.com/3_application/private/inventory.foocorp.com.key.pem \
--from-file=cert=inventory.foocorp.com/3_application/certs/inventory.foocorp.com.cert.pem \
--from-file=cacert=inventory.foocorp.com/2_intermediate/certs/ca-chain.cert.pem
```

これで、`frontend` と `inventory` をIstioリソースで公開する準備が整いました。まず、HTTPSトラフィック用にポート `443` の穴をあけるGatewayリソースを作成します。`mode: MUTUAL` は、内向きトラフィックに相互TLSを強制することを表すため注意して下さい。また、それぞれのサービスのために作成したSecretに対応する2つの異なる証明書のセットを指定します。

```YAML
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: foocorp-gateway
  namespace: default
spec:
  selector:
    istio: ingressgateway # use istio default ingress gateway
  servers:
  - port:
      number: 443
      name: https-frontend
      protocol: HTTPS
    tls:
      mode: MUTUAL
      credentialName: "frontend-credential"
    hosts:
    - "frontend.foocorp.com"
  - port:
      number: 443
      name: https-inventory
      protocol: HTTPS
    tls:
      mode: MUTUAL
      credentialName: "inventory-credential"
    hosts:
    - "inventory.foocorp.com"
```

次に、ゲートウェイからのルーティングを処理する2つのIstioのVirtualServiceを作成します。両方のサービスがゲートウェイのポート `443` にマッピングされているため、`host` ヘッダー（またはドメイン名）を使用して、要求するバックエンドサービスを指定します。

```YAML
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: frontend
spec:
  hosts:
  - "frontend.foocorp.com"
  gateways:
  - foocorp-gateway
  http:
  - match:
    - uri:
        exact: /
    route:
    - destination:
        host: frontend.ux.svc.cluster.local
        port:
          number: 80
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: inventory
spec:
  hosts:
  - "inventory.foocorp.com"
  gateways:
  - foocorp-gateway
  http:
  - match:
    - uri:
        exact: /
    route:
    - destination:
        host: inventory.corp-services.svc.cluster.local
        port:
          number: 80
```

これらのYAMLリソースを適用してから、`istio-ingressgateway` Pod内にある`ingress-sds` コンテナーのログを取得します。特定の証明書を使用してリソースを適用すると、SDSエージェントがそれらの証明書と鍵をIngressプロキシにマウントしたことがわかります。:

```bash
istio-ingressgateway-6f7d65d984-m2zmn ...
RESOURCE NAME:inventory-credential, EVENT: pushed key/cert pair to proxy
```

これで、クラスターの外部から2つのサービスにリクエストを送信する準備ができました。相互TLSを構成したので、サーバー（Ingress Gateway）がクライアントのIDを検証するために、`CA証明書` に加えて `証明書` と `鍵`を指定する必要があることに注意してください。

まず、クラスター外のホストから、frontendクライアント鍵を使用してfrontendをcurlでリクエストを送信します。：

```
$ curl -HHost:frontend.foocorp.com \
--resolve frontend.foocorp.com:$SECURE_INGRESS_PORT:$INGRESS_HOST \
--cacert frontend.foocorp.com/2_intermediate/certs/ca-chain.cert.pem \
--cert frontend.foocorp.com/4_client/certs/frontend.foocorp.com.cert.pem \
--key frontend.foocorp.com/4_client/private/frontend.foocorp.com.key.pem \
https://frontend.foocorp.com:$SECURE_INGRESS_PORT/

🏗 Welcome to FooCorp - Public Site
```

そして、内部のinventoryにinventoryのクライアント鍵を用いてアクセスします。：

```
$ curl -HHost:inventory.foocorp.com \
--resolve inventory.foocorp.com:$SECURE_INGRESS_PORT:$INGRESS_HOST \
--cacert inventory.foocorp.com/2_intermediate/certs/ca-chain.cert.pem \
--cert inventory.foocorp.com/4_client/certs/inventory.foocorp.com.cert.pem \
--key inventory.foocorp.com/4_client/private/inventory.foocorp.com.key.pem \
https://inventory.foocorp.com:$SECURE_INGRESS_PORT/

🛠 FooCorp - Inventory [INTERNAL]
```

ここで実際に何が起こっているのでしょうか？inventoryサービスを見て、Ingress Gatewayがクライアントを認証する方法を正確に見ていきましょう。

![](/images/secure-ingress-auth-steps.png)

1. クライアントはホスト `https://inventory.foocorp.com:443` にリクエストを送信します。
2. DNSは `inventory.foocorp.com` を Istio Ingress Gateway のパブリックIPアドレスに名前解決します（Istio Ingress Gateway はデフォルトでKubernetes Service の `type = LoadBalancer` としてプロビジョニングされています）。Ingress Gatewayは証明書と鍵をクライアントに提示します。
3. クライアントは、Ingress GatewayのIDを認証局（CA）で検証します。
4. クライアントは、証明書と鍵をIngress Gatewayに提示します。
5. サーバー（Ingress Gateway）は、クライアントのIDをCAで検証します。
6. クライアントとIngress Gatewayの間で安全な接続が確立され、Ingress Gatewayがリクエストを `inventory` サービスに転送します。

🎊できました！今後、新しいサービスを追加し続けることが可能となり、それに合わせて Ingress Gateway レプリカをスケールアウトできます。ご自身のクラスタは、セキュアで一元管理されたIngressを利用できます。

**詳しく学ぶ：**

- [Istio Ingress Gateway - コンセプト](https://istio.io/docs/concepts/traffic-management/#gateways)
- [Istio SDS Ingress、サーバー側TLSのみ](https://istio.io/docs/tasks/traffic-management/ingress/secure-ingress-sds/#configure-a-tls-ingress-gateway-for-multiple-hosts)
- [Istio SDS Ingress、手動ファイルマウントアプローチ](https://istio.io/docs/tasks/traffic-management/ingress/secure-ingress-mount/#before-you-begin)