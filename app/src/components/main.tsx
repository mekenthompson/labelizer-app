import * as React from "react";
import { Switch, Route, RouteProps, Redirect } from "react-router-dom"
import { Home } from "./home"
import { Setup } from "./setup"
import { Repos } from "./repos"
import { NotFound } from "./notFound"
import { Token } from "../auth/token"

export class Main extends React.Component<undefined, undefined> {
    shouldRedirect: boolean

    constructor(props: any){
        super(props);
        this.shouldRedirect = Token.storeFromQueryString()
    }

    render() {
        return (
            <main>
                {this.shouldRedirect ? (<Redirect to="/"/> ) : null}
                <Switch>
                    <Route exact path='/' component={Home}/>
                    <Route path='/setup' component={Setup}/>
                    <Route path='/repos' component={Repos} />
                    <Route path='*' component={NotFound} />
                </Switch>
            </main>
        )
    }
}