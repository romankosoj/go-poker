import { Container } from "pixi.js";
import { Player } from "./player";

class Players extends Container {
    constructor(state, id, table) {
        super();
        this.state = state;
        this.table = table;
        this.id = id;
        this.angles = [];
    }

    updateFromState() {
        if (this.state.state.players.length === 0) {
            return
        }
        let n = this.state.state.players.length
        let a = 360 / n;
        this.angles = [];
        this.removeChildren();
        for (let i = 0; i < n; i++) {
            this.angles.push(a * i * Math.PI / 180);
            let player = new Player(
                this.id,
                { angle: this.angles[i] },
                this.state,
                i
            );
            const x = this.table.x + (this.table.width + player.width) * Math.cos(this.angles[i]);
            const y = this.table.y + (this.table.height + player.height) * Math.sin(this.angles[i]);
            player.position.set(x, y);
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