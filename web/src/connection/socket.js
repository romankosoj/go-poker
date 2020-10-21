import { JOIN } from "../events/constants"
import { CreateEvent } from "../events/event";

class Game {
    constructor(state) {

        this.ws = new WebSocket("ws://localhost:8080/join")
        ws.onclose = e => {
            console.log("close", e)
        };

        ws.onopen = e => {
            console.log("open", e)

            if (e.type === "error") {
                return
            }

            ws.send(CreateEvent(JOIN, { username: "test", id: "testId" }))
        };

        ws.onmessage = e => {
            console.log("open", e)
            this.state.decodeChange(e);

        };

        ws.onerror = e => {
            console.log("open", e)
        };
    }
}

export { Game }