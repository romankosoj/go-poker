

export const

class GameEvent {
    constructor(name, data) {
        this.name = name;
        this.data = data;
    }

    toJSON() {
        return JSON.stringify({ event: this.name, data: this.data })
    }

    fromJSON(raw) {
        const event = JSON.parse(raw);
        this.name = event.name;
        this.data = event.data;
    }
}

export { GameEvent };
