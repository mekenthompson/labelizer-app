import * as React from "react";

export interface SetupProps { compiler: string; framework: string; }

export class Setup extends React.Component<SetupProps, undefined> {
    render() {
        return <h1>Hello from Setup!</h1>;
    }
}