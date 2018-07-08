import React, { Component } from "react";
import { connect } from "react-redux";
import { makeReservation } from "./store/reserve";
import { dateFromSlot } from "./store/funcs";
import store from './store';

const displayTime = (slot) => {
	const hour = Math.floor(slot.StartMinute / 60);
	const minute = slot.StartMinute % 60;
	return hour.toString().padStart(2, '0') + ':' + minute.toString().padStart(2, '0');
};

class ConnectedDay extends Component {
	constructor() {
		super();

		this.handleReserve = this.handleReserve.bind(this);
	}

	componentDidUpdate() {
		const end = performance.now();
		console.log('reservation took ' + (end - this.start) + 'ms');
	}

	handleReserve(e, slot) {
		this.start = performance.now();
		this.props.makeReservation(dateFromSlot(slot));
	}

	renderHours() {
		const now = new Date();
		const arr = [];
		for (let i=0; i < this.props.slots.length - 1; i += 2) {
			arr.push(
				<div key={i} className="hour">
					{this.renderSlot(now, this.props.slots[i])}
					{this.renderSlot(now, this.props.slots[i+1])}
				</div>
			)
		}
		return arr;
	}

	renderSlot(now, slot) {
		const date = dateFromSlot(this.props.slots[0]);
		if (now > date) {
			return <div className="reserved past">00:00</div>;
		}
		if (slot.Available) {
			return <button onClick={(e) => this.handleReserve(e, slot)}>{displayTime(slot)}</button>;
		} 
		return <div className={(slot.ClientID === this.props.clientID) ? "mine" : "reserved"}>{displayTime(slot)}</div>;
	}

	render() {
		const date = dateFromSlot(this.props.slots[0]);
		return <div className="day">
			<h4>{date.toLocaleDateString("en-us", {weekday: 'short', month: 'short', day: 'numeric'})}</h4>
			{this.renderHours()}
		</div>;
	}
}

const mapStateToProps = state => {
	return {
		clientID: state.ClientID
	};
};

const mapDispatchToProps = dispatch => {
	return {
		makeReservation: date => dispatch(makeReservation(store.getState(), date))
	};
};

const Day = connect(mapStateToProps, mapDispatchToProps)(ConnectedDay);

export default Day;