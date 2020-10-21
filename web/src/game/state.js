const { GAME_START, DEALER_SET, BLIND_SET, BIG_BLIND_SET, SMALL_BLIND_SET, WAIT_FOR_SMALL_BLIND_SET, WAIT_FOR_BIG_BLIND_SET, WAIT_FOR_PLAYER_ACTION, ACTION_PROCESSED, PLAYER_LEAVES, FLOP, TURN, RIVER, GAME_END } = require("../events/constants");
const { CHECK, BET, RAISE, FOLD } = require("../models/action");
const { Player } = require("../models/player");

export const UpdateEvents = {
    gameStart: 1,
    playerList: 2,
    player: 3,
    board: 4,
    gameEnd: 5,

}

export const GAMESTART = 0;
export const PLAYER_LIST = 1;
export

    class GameState {
    constructor(onUpdate) {
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
        this.onUpdate = onUpdate.bind(this)
        this.actionReceived = (action) => { };
        this.onGameEnd = () => {

        }
    }

    decodeChange(e) {
        if (e.name == GAME_START) {
            this.state.players = e.data.players;
            this.state.player = e.data.position;
        }
        switch (e.name) {
            case GAME_START:
                this.state.roundState = -1;
                this.state.player = e.data.position;
                for (let i = 0; i < e.data.players) {
                    this.state.players.push(new Player(e.data.players[i].username, e.data.players[i].id, 10));
                }
                this.onUpdate();
                break;
            case DEALER_SET:
                this.state.dealer = e.data;
                this.onUpdate();
                break;

            case SMALL_BLIND_SET:
                this.state.smallBlind = e.data;
                this.onUpdate();
                break;

            case BIG_BLIND_SET:
                this.bigBlind = e.data;
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
                this.state.players[e.position].waiting = true;
                this.onUpdate();
                break;

            case PLAYER_LEAVES:
                if (e.data.index === this.state.players.length - 1) {
                    this.state.players.pop();
                } else {
                    const end = this.state.players.slice(e.data.index, this.state.players.length - 1);
                    this.state.players.splice(i, 1, end);
                }
                this.onUpdate();

            case ACTION_PROCESSED:

                const lastIndex = this.state.lastAction
                this.state.players[lastIndex].isLastAction = false;
                this.onUpdate(UpdateEvents.player, lastAction);

                const i = this.state.waitingFor;
                const player = this.state.players[i];
                const action = new Action(e.data.action, e.data.position, e.data.amount);
                this.actionReceived(action);
                player.lastAction = action;
                player.isLastAction = true;
                if (action.action == FOLD) {
                    player.In = false;
                }
                if (action.action === BET || action.action === RAISE) {
                    this.state.players[i].bet = action.amount;
                    this.state.bet = action.amount;
                }

                this.state.lastAction = i;
                this.state.waitingFor = -1;
                this.onUpdate(UpdateEvents.player, i);
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
        }
    }
}