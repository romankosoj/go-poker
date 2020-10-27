import React from "react"

const { PLAYER_ACTION } = require("./events/constants")

class Action extends React.Component {
    constructor(props) {
        super(props);
        this.game = props.game
        this.state = {
            raise: 0,

        };
        this.foldPossible = false;
        this.checkPossible = false;
        this.callPossible = false;
        this.raisePossible = false;
        this.handleChange = this.handleChange.bind(this);
        this.fold = this.fold.bind(this);
        this.call = this.call.bind(this);
        this.check = this.check.bind(this);
        this.raise = this.raise.bind(this);

    }

    componentDidMount() {
        this.decodePossibleActions(this.props.actions)
    }

    decodePossibleActions(actions) {
        this.foldPossible = false;
        this.checkPossible = false;
        this.callPossible = false;
        this.raisePossible = false;
        console.log("actions: ", actions)

        if (actions & 1) {
            // 1 in first bit => Player can fold

            console.log("can fold", actions)

            this.foldPossible = true;
        }
        if ((actions >> 1) & 1) {
            // 1 in second bit => player can bet or call
            console.log("can call", actions)


            this.callPossible = true;
        }
        if ((actions >> 2) & 1) {

            console.log("can raise", actions)

            // 1 in second bit => player can bet or call
            this.raisePossible = true;
        }
        if ((actions >> 3) & 1) {
            console.log("can ceck", actions)

            // 1 in second bit => player can bet or call
            this.checkPossible = true;
        }
    }

    fold() {
        this.game.send({
            event: PLAYER_ACTION, data: {
                action: 1,
                payload: 0
            }
        })
    }

    call() {
        this.game.send({
            event: PLAYER_ACTION, data: {
                action: 2,
                payload: 0,
            }
        })
    }

    check() {
        this.game.send({
            event: PLAYER_ACTION, data: {
                action: 4,
                payload: 0,
            }
        })
    }

    raise() {
        this.game.send({
            event: PLAYER_ACTION, data: {
                action: 3,
                payload: this.state.raise,
            }
        })
    }

    handleChange(e) {
        this.setState({ raise: parseInt(e.target.value) });
    }

    render() {
        this.decodePossibleActions(this.props.actions);
        return (
            <div>
                { this.foldPossible ? <button onClick={this.fold}>FOLD</button> : <></>}
                { this.checkPossible ? <button onClick={this.check}>CHECK</button> : <></>}
                { this.callPossible ? <button onClick={this.call}>CALL</button> : <></>}
                { this.raisePossible
                    ? <div style={{ display: "inline" }}>
                        <input type="number" id="raiseInput" name="raiseInput" value={this.state.raise} onChange={this.handleChange} />
                        <button onClick={this.raise}>RAISE</button>
                    </div>
                    : <></>}
            </div>
        );
    }
}

export default Action