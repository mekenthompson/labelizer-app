import * as React from "react";
import { Link, Redirect } from "react-router-dom"
import { Navbar, Nav } from "react-bootstrap"
import { Token } from "../auth/token"

export interface HeaderProps { }

export class Header extends React.Component<HeaderProps, undefined> {
    handleLogout() {
        Token.clear()
    }

    render() {
        return (
            <Navbar inverse collapseOnSelect fixedTop>
                <Navbar.Header>
                    <Navbar.Brand>
                        <Link className="nav-link" to="/">Labelizer</Link>
                    </Navbar.Brand>
                    <Navbar.Toggle />
                </Navbar.Header>
                <Navbar.Collapse>
                    <Nav pullRight>
                        <li>
                            <a className="nav-link" href="https://github.com/apps/labelizer">Install Labelizer</a>
                        </li>
                        <li>
                            { !Token.isAuthenticated() ? 
                                (<a className="nav-link" href="/auth/github/signin">SignIn</a>) : 
                                (<a className="nav-link" href="/" onClick={this.handleLogout}>Logout</a>)}
                        </li>
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        )
    }
}