class Player {
    constructor(username, id, buyIn) {
        this.username = username
        this.id = id;
        this.buyIn = 0;
        this.bet = 0;
        this.cards = []
        this.in = true;
        this.lastAction;
        this.isLastAction = false;
        this.waiting = false;
    }

    setCards(cards) {
        this.cards = [];
        this.cards.concat(cards);
    }

    setBet(bet) {
        this.bet = bet;
    }

    setIn(val) {
        this.in = val;
    }


}

export { Player }