
const { JOIN_SUCCESS, GAME_START, DEALER_SET, WAIT_FOR_PLAYER_ACTION, ACTION_PROCESSED, FLOP, TURN, RIVER, GAME_END, HOLE_CARDS, PLAYER_JOIN } = require("../events/constants");
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
        this.name = "GameState test";
        this.onNotification = (text, st) => {

        }
        this.stateBuild = false;

        this.updateQueue = [];
    }

    setOnPossibleActions(onPossibleActions) {
        console.log("OnPossibleActions set");
        this.onPossibleAction = onPossibleActions.bind(this);
    }

    setOnGameStart(onGameStart) {
        this.onGameStart = onGameStart.bind(this);
    }

    setOnGameEnd(onGameEnd) {
        this.onGameEnd = onGameEnd.bind(this);
    }

    setOnUpdate(onUpdate) {
        this.onUpdate = onUpdate.bind(this);
    }

    setOnNotification(onNotification) {
        this.onNotification = onNotification;
    }

    decodeChange(e) {

        console.log("decoding event", e)
        switch (e.event) {

            case JOIN_SUCCESS:
                console.log("Players: ", e.data.players)
                this.state.roundState = -2;
                this.state.player = e.data.position;
                for (let i = 0; i < e.data.players.length; i++) {
                    this.state.players.push(new Player(e.data.players[i].username, e.data.players[i].id, e.data.players[i].buyIn));
                }
                this.onNotification("The game starts soon...", true);
                this.stateBuild = true;
                this.updateQueue.push({ event: UpdateEvents.lobbyJoin })
                break;

            case PLAYER_JOIN:
                this.state.players.push(new Player(e.data.player.username, e.data.player.id, e.data.player.buyIn))
                if (this.updateQueue[this.updateQueue.length - 1]?.event !== UpdateEvents.playerList) {
                    this.updateQueue.push({ event: UpdateEvents.playerList })
                }
                break;

            case GAME_START:
                this.state.roundState = -1;
                this.state.player = e.data.position;
                for (let i = 0; i < e.data.players.length; i++) {
                    this.state.players.push(new Player(e.data.players[i].username, e.data.players[i].id, e.data.players[i].buyIn));
                }
                this.onGameStart();
                break;
            case DEALER_SET:
                this.state.dealer = e.data;
                this.updateQueue.push({ event: UpdateEvents.dealer, data: e.data });
                break;

            case HOLE_CARDS:
                for (let i = 0; i < this.state.players.length; i++) {
                    this.state.players[i].cards = [{ color: -1, value: -1 }, { color: -1, value: -1 },]
                }
                this.state.players[this.state.player].cards = e.data.cards;
                this.updateQueue.push({ event: UpdateEvents.playerCards });
                break;

            case WAIT_FOR_PLAYER_ACTION:
                this.state.waitingFor = e.data.position;
                this.state.players[e.data.position].waiting = true;
                console.log("On possible actions called with", e.data.possibleActions)
                if (e.data.position === this.state.player) {
                    this.onPossibleAction(e.data.possibleActions);
                } else {
                    this.onPossibleAction(0);
                }
                this.updateQueue.push({ event: UpdateEvents.player, data: e.data.position })
                break;

            // case PLAYER_LEAVES:
            //     if (e.data.index === this.state.players.length - 1) {
            //         this.state.players.pop();
            //     } else {
            //         const end = this.state.players.slice(e.data.index, this.state.players.length - 1);
            //         this.state.players.splice(w, 1, end);
            //     }
            //     this.onUpdate(UpdateEvents.playerList);
            //     this.updateQueue.push(UpdateEvents.playerList)
            //     break;

            case ACTION_PROCESSED:
                if (this.state.lastAction > -1) {
                    const lastIndex = this.state.lastAction
                    this.state.players[lastIndex].isLastAction = false;
                    this.updateQueue.push({event: UpdateEvents.player, data: lastIndex});
                }
                const player = this.state.players[e.data.position];
                player.waiting = false;
                const action = new Action(e.data.action, e.data.position, e.data.amount);
                player.lastAction = action;
                player.isLastAction = true;
                if (action.action === FOLD) {
                    player.in = false;
                }
                if (action.action === BET || action.action === RAISE) {
                    this.state.players[e.data.position].bet = action.amount;
                    this.state.bet = action.amount;
                }

                this.state.lastAction = e.data.position;
                this.state.waitingFor = -1;
                this.updateQueue.push({ event: UpdateEvents.player, data: e.data.position });
                break;

            case FLOP:
                this.state.roundState = 1;
                this.state.board = e.data.cards;
                this.updateQueue.push({ event: UpdateEvents.board });

                break;
            case TURN:
                this.state.roundState = 2;
                this.state.board = this.state.board.concat(e.data.cards);
                this.updateQueue.push({ event: UpdateEvents.board });

                break;
            case RIVER:
                this.state.roundState = 3;
                this.state.board = this.state.board.concat(e.data.cards);
                this.updateQueue.push({ event: UpdateEvents.board });

                break;

            case GAME_END:
                this.state.roundState = 4;
                console.log("notification request")
                this.onNotification("Game ended. Next game coninues now.", false);
                this.updateQueue.push({ event: UpdateEvents.gameEnd })
                this.onGameEnd();
                break;
            default:
                break;
        }
    }

    getPlayerState(position) {
        return this.state.players[position];
    }
}

export { GameState }

export const UpdateEvents = {
    gameStart: 1,
    playerList: 2,
    player: 3,
    board: 4,
    gameEnd: 5,
    dealer: 6,
    playerCards: 7,
    lobbyJoin: 8,
}