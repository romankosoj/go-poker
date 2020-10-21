class Board {
    constructor(cards, marginX = 20, marginY = 20) {
        this.cards = [];
        for (let i = 0; i < cards.length; i++) {
            const card = new Card(cards[i].value, cards[i].color);
            card.setup();
            this.cards.push(card)
        }

        this.marginX = marginX;
        this.marginCards = 10;
        this.w = (cardWidth * this.cards.length) + (2 * this.marginX) + (this.cards.length * 2 * this.marginCards);
        this.h = cardHeight + 2 * marginY;
    }

    setup(x, y) {
        this.x = x;
        this.y = y;
    }

    draw() {
        if (this.x && this.y) {
            fill(10, 10, 10, 10);
            rect(this.x, this.y, this.w, this.h, 10);

            let itemX = this.x + this.marginX;
            let itemY = this.y + 20
            for (let i = 0; i < this.cards.length; i++) {
                itemX += this.marginCards;
                this.cards[i].draw(itemX, itemY);
                itemX += this.marginCards + cardWidth;
            }
        }
    }
}