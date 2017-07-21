import * as React from "react";
import { Switch, Route } from "react-router-dom"
import { Home } from "./home"
import { Setup } from "./setup"

export interface MainProps { }

export class Main extends React.Component<MainProps, undefined> {
    render() {
        return (
            <main>
                <Switch>
                    <Route exact path='/' component={Home}/>
                    <Route path='/setup' component={Setup}/>
                </Switch>
            </main>
        )
    }
}