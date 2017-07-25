import * as React from "react";

export class NotFound extends React.Component<undefined, undefined> {
    render() {
        return <div>
            <h1>404 - Page Not Found</h1>
            <p>I'm sorry, the page you were looking for cannot be found!</p>
        </div>
    }
}