import * as PIXI from "pixi.js"
import { Board } from "./board";
import { Notification } from "./notification";
import { Player } from "./player";
import { rW, rH, registerApp, isMobile } from "./utils";


let app = new PIXI.Application({
    resolution: 1,
    antialias: true,
    transparent: false,
    resizeTo: window,
});

app.renderer.backgroundColor = 0xffffff;
document.body.appendChild(app.view);

app.res

app.loader
    .add("textures/cards.json")
    .load(setup)

PIXI.Renderer.registerPlugin("interaction", PIXI.InteractionManager);

registerApp(app);

let angles = [];
let players = [];
let state;
function setup() {
    let id = app.loader.resources["textures/cards.json"].textures;

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
            username: "test",
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
    const playerIndex = 4;

    let table = new PIXI.Graphics();
    table.beginFill(0x1daf08);
    let tableWidth = 375;
    let tableHeight = 225;

    if (isMobile()) {
        tableWidth = 85;
        tableHeight = 50;
    }

    table.drawEllipse(tableWidth, tableHeight, tableWidth, tableHeight)
    //table.drawRoundedRect(0, 0, tableWidth * 2, tableHeight * 2, tableHeight)

    table.position.set(app.renderer.width / 2 - tableWidth, app.renderer.height / 2 - tableHeight)
    table.endFill();
    app.stage.addChild(table);

    let board = new Board(id, {});
    board.addCards(3);
    board.update({
        updatedWidth: () => {
            board.position.set((app.renderer.width / 2) - (board.width / 2), (app.renderer.height / 2) - (board.height / 2));
        }
    })
    board.position.set((app.renderer.width / 2) - (board.width / 2), (app.renderer.height / 2) - (board.height / 2));


    app.stage.addChild(board);

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
        generatePlayers(id, playersState, tableWidth - rW(75), tableHeight - rH(65), table.x + tableWidth, table.y + tableHeight);
    } else {
        generatePlayers(id, playersState, tableWidth, tableHeight, table.x + tableWidth, table.y + tableHeight);
    }


    let notification = new Notification(app.renderer.width, app.renderer.height)

    app.stage.addChild(notification)

    notification.position.set(0, 0)

    players[0].update({ loading: true })
    state = play;
    app.ticker.add(delta => gameLoop(delta))
}

const generatePlayers = (id, playersState, tWidth, tHeight, tX, tY) => {
    let n = playersState.length
    let a = 360 / n;
    for (let i = 0; i < n; i++) {
        angles.push(a * i * Math.PI / 180);
        let player = new Player(id, playersState[i], angles[i]);
        const x = tX + (tWidth + player.width) * Math.cos(angles[i]);
        const y = tY + (tHeight + player.height) * Math.sin(angles[i]);
        player.position.set(x, y);
        players.push(player);
        app.stage.addChild(player);
    }
}

function gameLoop(delta) {
    players[0].gameLoop(delta)
}

function play(delta) {

}
function end() {

}

