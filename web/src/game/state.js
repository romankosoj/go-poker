const { GAME_START, DEALER_SET, BIG_BLIND_SET, SMALL_BLIND_SET, WAIT_FOR_SMALL_BLIND_SET, WAIT_FOR_BIG_BLIND_SET, WAIT_FOR_PLAYER_ACTION, ACTION_PROCESSED, PLAYER_LEAVES, FLOP, TURN, RIVER, GAME_END, HOLE_CARDS } = require("../events/constants");
const { BET, RAISE, FOLD, Action } = require("../models/action");
const { Player } = require("../models/player");

class GameState {
    constructor() {
        this.state = {
            players: [],
            player: -1,
            dealer: -1,
            board: [],
            bet: 0,
            roundState: 0,
            waitingFor: -1,
            smallBlind: 0,
            bigBlind: 0,
            lastAction: -1,
        };
        this.onPossibleAction = (actions) => {};
        this.name = "GameState test";
    }

    setOnGameStart(onGameStart){
        this.onGameStart = onGameStart.bind(this);
    }

    setOnGameEnd(onGameEnd){
        this.onGameEnd = onGameEnd.bind(this);
    }

    setOnUpdate(onUpdate){
        this.onUpdate = onUpdate.bind(this);
    }

    decodeChange(e) {

        console.log("decoding event", e)
        const w = this.state.waitingFor;
        console.log(e.event);
        switch (e.event) {
            case GAME_START:
                this.state.roundState = -1;
                this.state.player = e.data.position;
                for (let i = 0; i < e.data.players.length; i++) {
                    this.state.players.push(new Player(e.data.players[i].username, e.data.players[i].id, 10));
                }
                this.onGameStart();
                setTimeout(() => {
                    this.onUpdate(UpdateEvents.gameStart);
                }, 500);
                break;
            case DEALER_SET:
                this.state.dealer = e.data;
                this.onUpdate(UpdateEvents.dealer, e.data);
                break;

            case SMALL_BLIND_SET:
                this.state.smallBlind = e.data;
                this.state.players[w].bet = e.data;
                this.state.players[w].lastAction = new Action(BET, w, e.data);
                this.state.waitingFor = -1;
                this.onUpdate(UpdateEvents.player, w);
                break;

            case BIG_BLIND_SET:
                this.bigBlind = e.data;
                this.state.bigBlind = e.data;
                this.state.players[w].bet = e.data;
                this.state.players[w].lastAction = new Action(RAISE, w, e.data);
                this.state.waitingFor = -1;
                this.onUpdate();
                break;

            case WAIT_FOR_SMALL_BLIND_SET:
                this.state.waitingFor = e.data;
                this.state.players[e.data].waiting = true;
                this.onUpdate(UpdateEvents.player, e.data);
                break;

            case WAIT_FOR_BIG_BLIND_SET:
                this.state.waitingFor = e.data;
                this.state.players[e.data].waiting = true;
                this.onUpdate(UpdateEvents.player, e.data);
                break;

            case HOLE_CARDS:
                this.state.players[this.state.player].cards = e.data;
                this.onUpdate(UpdateEvents.player, this.state.player);
                break;

            case WAIT_FOR_PLAYER_ACTION:
                this.state.waitingFor = e.data.position;
                this.state.players[e.data.position].waiting = true;
                this.onPossibleAction(e.data.possibleActions);
                this.onUpdate(UpdateEvents.player, e.data.position);
                break;

            case PLAYER_LEAVES:
                if (e.data.index === this.state.players.length - 1) {
                    this.state.players.pop();
                } else {
                    const end = this.state.players.slice(e.data.index, this.state.players.length - 1);
                    this.state.players.splice(w, 1, end);
                }
                this.onUpdate();
                break;

            case ACTION_PROCESSED:
                if (this.state.lastAction > -1) {
                    const lastIndex = this.state.lastAction
                    this.state.players[lastIndex].isLastAction = false;
                    this.onUpdate(UpdateEvents.player, lastIndex);
                }
                const player = this.state.players[e.data.position];
                player.waiting = false;
                const action = new Action(e.data.action, e.data.position, e.data.amount);
                player.lastAction = action;
                player.isLastAction = true;
                if (action.action === FOLD) {
                    player.In = false;
                }
                if (action.action === BET || action.action === RAISE) {
                    this.state.players[e.data.position].bet = action.amount;
                    this.state.bet = action.amount;
                }

                this.state.lastAction = e.data.position;
                this.state.waitingFor = -1;
                this.onUpdate(UpdateEvents.player, e.data.position);
                break;

            case FLOP:
                this.state.roundState = 1;
                this.state.board.concat(e.data.cards);
                this.onUpdate(UpdateEvents.board);
                break;
            case TURN:
                this.state.roundState = 2;
                this.state.board.concat(e.data.cards);
                this.onUpdate(UpdateEvents.board);
                break;
            case RIVER:
                this.state.roundState = 3;
                this.state.board.concat(e.data.cards);
                this.onUpdate(UpdateEvents.board);
                break;

            case GAME_END:
                this.state.roundState = 4;
                this.onUpdate(UpdateEvents.gameEnd);
                this.onGameEnd();
                break;
            default:
                break;
        }
    }

    getPlayerState(position){
        return this.state.players[position];
    }
}

export { GameState}

export const UpdateEvents = {
    gameStart: 1,
    playerList: 2,
    player: 3,
    board: 4,
    gameEnd: 5,
    dealer: 6,

}