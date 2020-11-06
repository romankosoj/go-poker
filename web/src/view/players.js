import { Container } from "pixi.js";
import { Player } from "./player";

class Players extends Container {
    constructor(state, table) {
        super();
        this.state = state;
        this.table = table;
        this.angles = [];
    }

    setup(id) {
        this.id = id;
    }

    updateFromState() {
        console.log("Update Players from State: ", this.state.state.players)

        if (this.state.state.players.length === 0) {
            return
        }
        const a = 360 / this.state.state.players.length;
        this.angles = [];
        this.removeChildren();
        for (let i = 0; i < this.state.state.players.length; i++) {
            this.angles.push(a * i * Math.PI / 180);
            let player = new Player(
                this.id,
                { angle: this.angles[i] },
                this.state,
                i
            );
            const xPos = this.table.x + (this.table.width + player.width) * Math.cos(this.angles[i]);
            const yPos = this.table.y + (this.table.height + player.height) * Math.sin(this.angles[i]);
            player.position.set(xPos, yPos);
            this.addChild(player);
        }
    }

    updatePlayerFromState(i) {
        if (i >= this.children.length) {
            return;
        }
        this.children[i].updateFromState();
    }

    updateAllPlayersFromState() {
        for (let i = 0; i < this.children.length; i++) {
            this.updatePlayerFromState(i)
        }
    }

    gameLoop(delta) {
        for (let i = 0; i < this.children.length; i++) {
            this.children[i].gameLoop();
        }
    }
}

export { Players }