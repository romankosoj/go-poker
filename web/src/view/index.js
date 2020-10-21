import { Application, Graphics, InteractionManager, Renderer } from "pixi.js";
import React from "react";
import { Board } from "./board";
import { Player } from "./player";
import { rW, rH, registerApp, isMobile } from "./utils";

class View extends React.Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        const d = document.getElementById("view");
        this.app = new Application({
            antialias: true,
            transparent: false,
            resolution: 1,
            resizeTo: window,
        });

        this.app.renderer.backgroundColor = 0xffffff;

        d.appendChild(this.app.view);

        this.app.loader
            .add("textures/cards.json")
            .load(this.setup.bind(this))

        Renderer.registerPlugin("interaction", InteractionManager);

        registerApp(this.app);

        this.angles = [];
        this.players = [];
    }

    setup() {
        let id = this.app.loader.resources["textures/cards.json"].textures;

        rH(0);

        let playersState = [
            {
                username: "test",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 5, color: 3 },
                ],
            },
            {
                username: "test",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 12, color: 3 },
                ],
            },
            {
                username: "test",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 12, color: 3 },
                ],
            },
            {
                username: "test2",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 12, color: 3 },
                ],
            }, {
                username: "you",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 12, color: 3 },
                ],
            }, {
                username: "test",
                bet: 120,
                cards: [
                    { value: 1, color: 2 },
                    { value: 12, color: 3 },
                ],
            },
        ];

        let table = new Graphics();
        table.beginFill(0x1daf08);
        const tableWidth = rW(375);
        const tableHeight = rH(225);
        table.drawEllipse(tableWidth, tableHeight, tableWidth, tableHeight)
        //table.drawRoundedRect(0, 0, tableWidth * 2, tableHeight * 2, tableHeight)
        table.position.set(this.app.renderer.width / 2 - tableWidth, this.app.renderer.height / 2 - tableHeight)
        table.endFill();
        this.app.stage.addChild(table);

        let board = new Board(id, {});
        board.addCards(3);
        board.update({
            updatedWidth: () => {
                board.position.set((this.app.renderer.width / 2) - (board.width / 2), (this.app.renderer.height / 2) - (board.height / 2));
            }
        })
        board.position.set((this.app.renderer.width / 2) - (board.width / 2), (this.app.renderer.height / 2) - (board.height / 2));


        this.app.stage.addChild(board);

        setTimeout(() => {
            board.pushOrUpdate([
                { value: 2, color: 1 },
                { value: 11, color: 3 },
                { value: 12, color: 2 },
                { value: 12, color: 2 },
            ])
        }, 2000)

        console.log(isMobile());

        if (isMobile()) {
            this.generatePlayers(id, playersState, tableWidth - rW(75), tableHeight - rH(65), table.x + tableWidth, table.y + tableHeight);
        } else {
            this.generatePlayers(id, playersState, tableWidth, tableHeight, table.x + tableWidth, table.y + tableHeight);
        }
        this.app.ticker.add(delta => this.gameLoop(delta))
    }

    generatePlayers(id, playersState, tWidth, tHeight, tX, tY) {
        let n = playersState.length
        let a = 360 / n;
        for (let i = 0; i < n; i++) {
            this.angles.push(a * i * Math.PI / 180);
            let player = new Player(id, playersState[i], this.angles[i]);
            const x = tX + (tWidth + player.width) * Math.cos(this.angles[i]);
            const y = tY + (tHeight + player.height) * Math.sin(this.angles[i]);
            player.position.set(x, y);
            this.players.push(player);
            this.app.stage.addChild(player);
        }
    }

    render() {
        return (
            <div id="view">
            </div>
        )
    }
}

export default View;