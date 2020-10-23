import { Container, Graphics, Text } from "pixi.js"
import { Board } from "./board";
import { isMobile } from "./utils";

class Player extends Container {
    constructor(id, options, state, index) {
        super();

        this.state = state;
        this.playerState = this.state.getPlayerState(index);
        this.index = index;

        this.options = {
            width: 250,
            marginX: 10,
            marginY: 8,
            avatarRadius: 20,
            angle: 0,
            fontFamily: "Source Sans Pro",
            fontSize: 20,
            foreground: 0x000000,
            fillAvatar: 0xffffff,
            waitingAnimationVelocity: 5,
            background: {
                value: 0x000000,
                alpha: 0.12,
            }
        }

        if (isMobile()) {
            this.options.width = 150;
        }

        this.betLabel = new Text("", { align: "center" });
        this.usernameLabel = new Text("", { align: "center" });
        this.avatar = new Graphics();
        this.board = new Board(id, {
            paddingX: 8,
            paddingY: 5,
            paddingCardX: 5,
        });
        this.background = new Graphics();
        this.waiting = new Graphics();
        this.waitingVelocity = this.options.waitingAnimationVelocity;
        this.waitingAnimeDir = true;

        this.topRow = new Container()
        this.topRow.addChild(this.avatar, this.usernameLabel, this.betLabel, this.waiting);
        this.addChild(this.background, this.topRow, this.board);

        this.update(options)
    }

    gameLoop(delta) {
        if (this.waitingVelocity < 3.25) {
            this.waitingAnimeDir = true;
        } else if (this.waitingVelocity > 12) {
            this.waitingAnimeDir = false
        }
        if (this.waitingAnimeDir) {
            this.waitingVelocity += 0.075
        } else {
            this.waitingVelocity -= 0.075
        }
        this.waiting.angle = (this.waiting.angle + this.waitingVelocity + delta) % 360
    }


    update(options) {
        this.options = {
            ...this.options,
            ...options,
        };

        this.calcWidth = this.options.width - 2 * this.options.marginX;

        this.usernameLabel.text = this.playerState.username
        this.usernameLabel.style = {
            fontFamily: this.options.fontFamily,
            fontSize: this.options.fontSize,
            fill: this.options.foreground,
        }
        this.usernameLabel.pivot.set(this.usernameLabel.width, 0);

        this.betLabel.text = this.playerState.bet
        this.betLabel.style = {
            fontFamily: this.options.fontFamily,
            fontSize: this.options.fontSize,
            fill: this.state.state.lastIndex === this.index ? 0xff0000 : this.options.foreground
        }
        this.betLabel.pivot.set(this.betLabel.width, 0);

        this.avatar.clear()
        this.avatar.beginFill(this.options.fillAvatar);
        this.avatar.drawCircle(this.options.avatarRadius, this.options.avatarRadius, this.options.avatarRadius);
        this.avatar.endFill();

        // updateCards
        this.board.pushOrUpdate(this.playerState.cards);


        this.board.position.set((this.calcWidth / 2) - (this.board.width / 2) + this.options.marginX, this.avatar.height + 2 * this.options.marginY);
        this.avatar.position.set(this.options.marginX, this.options.marginY);
        this.usernameLabel.position.set(this.calcWidth - this.options.avatarRadius, this.avatar.height / 2)
        this.betLabel.position.set(this.calcWidth - this.options.avatarRadius - this.usernameLabel.width, this.avatar.height / 2)

        if (this.playerState.waiting) {
            this.waiting.lineStyle(4, 0x000000);
            this.waiting.arc(this.options.avatarRadius, this.options.avatarRadius, this.options.avatarRadius, 0, Math.PI);
            this.waiting.position.set(this.avatar.position.x + this.options.avatarRadius, this.avatar.position.y + this.options.avatarRadius)
            this.waiting.pivot.set(this.options.avatarRadius)
        } else {
            this.waiting.clear();
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

    updateFromState() {
        this.playerState = this.state.getPlayerState(this.index);
        this.update({});
        console.log("Waiting for player ?:", this.playerState.waiting, " [", this.index, "]");
    }
}

export { Player }