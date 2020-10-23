import { Container, Graphics, Text } from "pixi.js"
import { Board } from "./board";

class Player extends Container {
    constructor(id, options, state, index) {
        super();
        
        this.state = state;
        this.playerState = this.state.getPlayerState(index);
        this.index = index;

        this.options = {
            marginX: 10,
            marginY: 8,
            avatarRadius: 20,
            angle: 0,
            fontFamily: "Source Sans Pro",
            fontSize: 20,
            foreground: 0x000000,
            fillAvatar: 0xffffff,
            background: {
                value: 0x000000,
                alpha: 0.12,
            }
        }

        this.usernameLabel = new Text("", { align: "center" });
        this.avatar = new Graphics();
        this.board = new Board(id, {
            paddingX: 8,
            paddingY: 5,
            paddingCardX: 5,
        });
        this.background = new Graphics();


        this.addChild(this.background, this.avatar, this.usernameLabel, this.board);

        this.update(options)
    }


    update(options) {
        this.options = {
            ...this.options,
            ...options,
        };

        this.left = this.options.angle > Math.PI / 2 && this.options.angle <= Math.PI + Math.PI / 2;

        this.usernameLabel.text = this.playerState.username
        this.usernameLabel.style = {
            fontFamily: this.options.fontFamily,
            fontSize: this.options.fontSize,
            fill: this.options.foreground,
        }
        this.avatar.clear()
        this.avatar.beginFill(this.options.fillAvatar);
        this.avatar.drawCircle(this.options.avatarRadius, this.options.avatarRadius, this.options.avatarRadius);
        this.avatar.endFill();
        
        // updateCards
        this.board.pushOrUpdate(this.options.cards);
        
        
        this.board.position.set(this.options.marginX, this.avatar.height + 2 * this.options.marginY);
        if (this.left) {
            this.usernameLabel.pivot.set(0,0);
            this.avatar.position.set(this.board.width - this.options.avatarRadius, this.options.marginY)
            this.usernameLabel.position.set(this.options.marginX, this.avatar.height / 2);
        } else {
            this.usernameLabel.pivot.set(this.usernameLabel.width,0);
            this.avatar.position.set(this.options.marginX, this.options.marginY);
            this.usernameLabel.position.set(this.board.width, this.avatar.height / 2)
        }

        this.background.clear();
        this.background.beginFill(this.options.background.value, this.options.background.alpha);
        this.background.drawRoundedRect(0, 0, this.board.width + 2 * this.options.marginX, this.avatar.height + 3 * this.options.marginY + this.board.height, 10);
        this.background.endFill();

        this.onResize();

    }

    onResize() {
        this.pivot.set(this.width * 0.5, this.height * 0.5);
    }

    updateFromState(){
        this.playerState = this.state.getPlayerState(this.index);
        this.update({});
        console.log("Waiting for player ?:", this.playerState.waiting, " [", this.index, "]");
    }
}

export { Player }