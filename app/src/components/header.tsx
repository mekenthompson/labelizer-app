import * as React from "react";
import { Link } from "react-router-dom"
import { Navbar, Nav, NavItem, NavDropdown } from "react-bootstrap"

export interface HeaderProps { }

export class Header extends React.Component<HeaderProps, undefined> {
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
                            <Link className="nav-link" to="/setup">Setup</Link>
                        </li>
                        <li>
                            <a className="nav-link" href="https://github.com/apps/labelizer">Install</a>
                        </li>
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        )
    }
}