import React from "react"
import { Link } from "gatsby"
import logo from "../../static/logo.png" // Tell Webpack this JS file uses this image

import { rhythm, scale } from "../utils/typography"

class Layout extends React.Component {
  render() {
    const { location, title, children } = this.props
    const rootPath = `${__PATH_PREFIX__}/`
    let header
    let footer

    if (location.pathname === rootPath) {
      header = (
        <div>
          <Link
            style={{
              boxShadow: `none`,
              textDecoration: `none`,
              color: `inherit`,
              marginBottom: 0,
              marginTop: 0
            }}
            to={`/`}
          >
          <img src={logo} alt="Logo" style={{
                      marginBottom: 0,
          }}/>
          </Link>
          </div>
      )
    } else {
      header = (
        <h3
          style={{
            fontFamily: `Noto Sans, sans-serif`,
            marginTop: 0,
          }}
        >
          <Link
            style={{
              boxShadow: `none`,
              textDecoration: `none`,
              color: `inherit`,
            }}
            to={`/`}
          >
               <img src={logo} alt="Logo" />
          </Link>
        </h3>
      )
    }

    if (location.pathname === rootPath) {
      footer = (
        <div></div>
      )
    } else {
    footer = (
      <p       style={{
        fontFamily: `Inconsolata, monospace`,
        fontSize: '1.0em',
        marginTop: 0,
      }}> Made with 💙 by <a href="https://twitter.com/askmeegs">Megan O'Keefe</a> — <a rel="license" href="http://creativecommons.org/licenses/by/3.0/">License</a>
</p>
    )
    }
    return (
      <div
        style={{
          marginLeft: `auto`,
          marginRight: `auto`,
          maxWidth: rhythm(24),
          padding: `${rhythm(1.5)} ${rhythm(3 / 4)}`,
        }}
      >
        <header>{header}</header>
        <main>{children}</main>
        <footer>{footer}</footer>
      </div>
    )
  }
}

export default Layout
