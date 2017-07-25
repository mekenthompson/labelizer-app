import * as QS from "query-string"

const TOKEN = "token"

export class Token {

    static storeFromQueryString(): boolean {
        let qs = QS.parse(window.location.search)
        let token = qs[TOKEN]
        if (token != null) {
            this.store(token)
            return true
        }
        return false
    }

    static store(token: string) {
        localStorage.setItem(TOKEN, token)
    }

    static fetch(): string {
        return localStorage.getItem(TOKEN)
    }

    static isAuthenticated(): boolean {
        return this.fetch() != null
    }

    static clear(): boolean {
        localStorage.removeItem(TOKEN)
        return true
    }
}