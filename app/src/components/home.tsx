import * as React from "react";

export interface HomeProps { compiler: string; framework: string; }

export class Home extends React.Component<HomeProps, undefined> {
    render() {
        return <h1>Hello from Home!</h1>;
    }
}