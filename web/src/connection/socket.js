import { JOIN } from "../events/constants"
import { SendCreateEvent, SendEvent } from "../events/event";

class Game {
    constructor(state, entry, onClose) {
        this.state = state;
        this.credentials = entry;
        this.onClose = onClose.bind(this);
        this.started = new Promise((resolve, reject) => {
            this.state.setOnGameStart(() => {
                console.log("Resolving game start");
                resolve();
            });
        });
    }

    start() {
        console.log(this.state.name);
        this.ws = new WebSocket("ws://localhost:8080/join")
        this.ws.onclose = e => {
            console.log("close", e)
        };

        this.ws.onmessage = e => {
            if (e.data) {
                this.state.decodeChange(JSON.parse(e.data));
            }
        };

        this.ws.onopen = e => {
            if (e.type === "error") {
                this.onClose();
                return
            }

            // Join event
            this.ws.send(SendCreateEvent(JOIN, this.credentials));
        };


        this.ws.onerror = e => {
            console.log("error", e)
            this.onClose();
        };
    }

    send(event) {
        this.ws.send(SendEvent(event));
    }
}

export { Game }