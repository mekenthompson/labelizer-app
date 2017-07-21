import * as React from "react";

import { Header } from "./header"
import { Main } from "./main"

export interface AppProps { }

export class App extends React.Component<AppProps, undefined> {
    render() {
        return (
        <div>
            <Header />
            <Main />
        </div>)
    }
}