import * as React from "react";
import { Link } from "react-router-dom"

export interface HeaderProps { }

export class Header extends React.Component<HeaderProps, undefined> {
    render() {
        return (
            <header>
                <nav>
                    <ul>
                        <li><Link to='/'>Home</Link></li>
                        <li><Link to='/setup'>Setup</Link></li>
                    </ul>
                </nav>
            </header>
        )
    }
}