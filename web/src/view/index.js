import { Application, Graphics, InteractionManager, Loader, Renderer } from "pixi.js";
import React from "react";
import { Board } from "./board";
import { Player } from "./player";
import { rW, rH, registerApp, isMobile } from "./utils";
import PropTypes from 'prop-types';
import { Game } from "../connection/socket";
import { UpdateEvents } from "../game/state";

class View extends React.Component {
    constructor(props) {
        super(props);
        console.log("Game", props.game);
        console.log("State", props.game.state);
        this.gameState = props.game.state;
        this.game = props.game.state.state;

        this.angles = [];
        this.players = [];
    }

    componentDidMount() {
        this.props.game.started.then(() => {
            console.log("game started in view");
            const d = document.getElementById("view");

            this.app = new Application({
                antialias: true,
                transparent: false,
                resolution: 1,
                resizeTo: window,
            });

            this.gameState.setOnUpdate(this.gameUpdate.bind(this));

            this.app.loader = this.props.loader;
            this.app.loader.load(this.setup.bind(this))

            Renderer.registerPlugin("interaction", InteractionManager);

            this.app.renderer.backgroundColor = 0xffffff;

            d.appendChild(this.app.view);

            registerApp(this.app);
        })
    }

    gameUpdate(event, data) {
        if (this.players.length < 0) {
            this.updatePlayers();
        }
        console.log("Game Update in view [", event, "]: ", data);
        if (event === UpdateEvents.playerList) {
            this.updatePlayers();
        }
        if (event === UpdateEvents.playerCards) {
            for (let i = 0; i < this.players.length; i++) {
                this.players[i].updateCardsFromState();
            }
        }
        if (event === UpdateEvents.player) {
            console.log(this.players);
            this.players[data].updateFromState();
        }
        if (event === UpdateEvents.dealer) {
            this.players[data].updateFromState();
        }
        if (event === UpdateEvents.board) {
            this.board.updateFromState();
        }
    }


    setup() {
        this.id = this.app.loader.resources["textures/cards.json"].textures;

        this.table = new Graphics();
        this.table.beginFill(0x1daf08);
        this.tableWidth = rW(375);
        this.tableHeight = rH(225);
        this.table.drawEllipse(this.tableWidth, this.tableHeight, this.tableWidth, this.tableHeight)
        //table.drawRoundedRect(0, 0, this.tableWidth * 2, this.tableHeight * 2, this.tableHeight)
        this.table.position.set(this.app.renderer.width / 2 - this.tableWidth, this.app.renderer.height / 2 - this.tableHeight)
        this.table.endFill();
        this.app.stage.addChild(this.table);

        this.board = new Board(this.id, {}, this.game);
        this.board.addCards(3);
        this.board.update({
            updatedWidth: () => {
                this.board.position.set((this.app.renderer.width / 2) - (this.board.width / 2), (this.app.renderer.height / 2) - (this.board.height / 2));
            }
        })
        this.board.position.set((this.app.renderer.width / 2) - (this.board.width / 2), (this.app.renderer.height / 2) - (this.board.height / 2));


        this.app.stage.addChild(this.board);

        this.updatePlayers();
        //this.app.ticker.add(delta => this.gameLoop(delta))
    }

    updatePlayers() {
        if (isMobile()) {
            this.generatePlayers(this.id, this.game.players, this.tableWidth - rW(75), this.tableHeight - rH(65), this.table.x + this.tableWidth, this.table.y + this.tableHeight)
        } else {
            this.generatePlayers(this.id, this.game.players, this.tableWidth, this.tableHeight, this.table.x + this.tableWidth, this.table.y + this.tableHeight)
        }
    }

    generatePlayers(id, players, tWidth, tHeight, tX, tY) {
        if (this.game.players && this.game.players.length) {
            this.players = [];
            this.angles = [];
            let n = players.length
            let a = 360 / n;
            for (let i = 0; i < n; i++) {
                this.angles.push(a * i * Math.PI / 180);
                let player = new Player(id,
                    {
                        angle: this.angles[i],
                    },
                    this.gameState,
                    i
                );
                const x = tX + (tWidth + player.width) * Math.cos(this.angles[i]);
                const y = tY + (tHeight + player.height) * Math.sin(this.angles[i]);
                player.position.set(x, y);
                this.players.push(player);
                this.app.stage.addChild(player);
            }
        }
    }

    render() {
        return (
            <div id="view">
            </div>
        )
    }
}

View.propTypes = {
    game: PropTypes.instanceOf(Game),
    loader: PropTypes.instanceOf(Loader),
}

export default View;