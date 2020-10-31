import { Container, Graphics, Text } from "pixi.js"
import { Ease, ease } from "pixi-ease"
import { Board } from "./board";
import { isMobile } from "./utils";

class Player extends Container {
    constructor(id, options) {
        super();

        this.options = {
            width: 250,
            marginX: 10,
            marginY: 8,
            avatarRadius: 20,
            angle: 0,
            fontFamily: "Source Sans Pro",
            fontSize: "1.5rem",
            foreground: 0x000000,
            fillAvatar: 0xffffff,
            background: {
                value: 0x000000,
                alpha: 0.12,
            },
        }

        if (isMobile()) {
            this.options.width = 150
        }

        this.betLabel = new Text("20", { align: "center" });
        this.usernameLabel = new Text("", { align: "center" });
        this.avatar = new Graphics();
        this.board = new Board(id, {
            paddingX: 8,
            paddingY: 5,
            paddingCardX: 5,
        });
        this.background = new Graphics();

        this.loading = new Graphics();
        this.loadingR = 5;
        this.loadingRR = true;

        this.dealerButton = new Graphics();

        this.topRow = new Container();


        this.topRow.addChild(this.avatar, this.usernameLabel, this.loading, this.betLabel, this.dealerButton)

        this.addChild(this.background, this.topRow, this.board);

        this.update(options)
    }


    gameLoop(delta) {
        if (this.loadingR < 3.25) {
            this.loadingRR = true;
        } else if (this.loadingR > 12) {
            this.loadingRR = false
        }
        if (this.loadingRR) {
            this.loadingR += 0.075
        } else {
            this.loadingR -= 0.075
        }
        this.loading.angle = (this.loading.angle + this.loadingR + delta) % 360
    }


    update(options) {
        this.options = {
            ...this.options,
            ...options,
        };

        this.calcWidth = this.options.width - this.options.marginX * 2;

        this.dealerButton.clear()
        this.dealerButton.lineStyle(1, 0xff0000, 1);
        this.dealerButton.beginFill(0xffff00)
        this.dealerButton.drawCircle(0, 0, 12)
        this.dealerButton.endFill();


        this.usernameLabel.text = this.options.username
        this.usernameLabel.style = {
            fontFamily: this.options.fontFamily,
            fontSize: this.options.fontSize,
            fill: this.options.foreground,
        }

        this.betLabel.text = "30"
        this.betLabel.style = {
            fontFamily: this.options.fontFamily,
            fontSize: this.options.fontSize,
            fill: 0xff0000,
        }
        this.betLabel.pivot.set(0, 0)

        this.avatar.clear()
        this.avatar.beginFill(this.options.fillAvatar);
        this.avatar.drawCircle(this.options.avatarRadius, this.options.avatarRadius, this.options.avatarRadius);
        this.avatar.endFill();

        // updateCards
        this.board.pushOrUpdate(this.options.cards);


        this.board.position.set((this.calcWidth / 2) - (this.board.width / 2) + this.options.marginX, this.options.avatarRadius + 3 * this.options.marginY);
        this.avatar.position.set(this.options.marginX, this.options.marginY);
        this.usernameLabel.position.set(this.calcWidth - this.options.avatarRadius, this.options.avatarRadius + this.options.marginY)
        this.betLabel.position.set(this.calcWidth - this.options.avatarRadius - this.usernameLabel.width, this.options.avatarRadius + this.options.marginY)

        this.dealerButton.position.set(this.calcWidth - this.options.avatarRadius - this.usernameLabel.width - this.betLabel.width - this.options.marginX, this.options.avatarRadius + this.options.marginY)

        if (this.options.loading) {
            this.loading.lineStyle(4, 0x000000);
            this.loading.arc(this.options.avatarRadius, this.options.avatarRadius, this.options.avatarRadius, 0, Math.PI);
            this.loading.position.set(this.avatar.position.x + this.options.avatarRadius, this.avatar.position.y + this.options.avatarRadius)
            this.loading.pivot.set(this.options.avatarRadius)
        } else {
            this.loading.clear();
        }

        this.background.clear();
        this.background.beginFill(this.options.background.value, this.options.background.alpha);
        this.background.drawRoundedRect(0, 0, this.options.width, this.topRow.height + 3 * this.options.marginY + this.board.height, 10);
        this.background.endFill();

        this.onResize();



    }

    onResize() {
        this.pivot.set(this.width * 0.5, this.height * 0.5);
    }
}

export { Player }