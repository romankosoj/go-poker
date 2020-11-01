import { Player } from "./player";

const { Container } = require("pixi.js");

class Players extends Container {
    constructor(id, table) {
        super();

        this.table = table;

        this.angles = [];

        this.id = id;

    }

    updatePlayers(players) {
        let n = players.length
        let a = 360 / n;
        this.angles = [];
        for (let i = 0; i < n; i++) {
            this.angles.push(a * i * Math.PI / 180);
            let player = new Player(this.id, players[i], this.angles[i]);
            const x = this.table.x + (this.table.width + player.width) * Math.cos(this.angles[i]);
            const y = this.table.y + (this.table.height + player.height) * Math.sin(this.angles[i]);
            player.position.set(x, y);
            this.addChild(player);
        }
    }

    gameLoop(delta) {
        for (let i = 0; i < this.children.length; i++) {
            this.children[i].gameLoop();
        }
    }

    updatePlayerIndex(i) {
        let player = this.children[i];
        player.update({ loading: true });
    }
}

const generatePlayers = (id, playersState, tWidth, tHeight, tX, tY) => {
}


export { Players }