import React, { Component } from "react";
import { connect } from "react-redux";
import Day from "./Day.jsx";
import { fetchWeek } from "./store/schedule";
import store from './store';
import { dateFromSlot } from "./store/funcs";

class ConnectedSchedule extends Component {
	constructor() {
		super();

		this.handlePaging = this.handlePaging.bind(this);
	}

	handlePaging(e) {
		this.props.fetchWeek(this.props.anchorDate, e.target.className);
	}

	renderWeek() {
		const arr = [];
		for (let i=0; i < this.props.days.length; ++i) {
			arr.push(<Day key={i} slots={this.props.days[i]} />);
		}
		return arr;
	}

	render() {
		return <div>
			<button className="prev" onClick={this.handlePaging}>Previous Week</button>
			<div id="week">
				{this.renderWeek()}
			</div>
			<button className="next" onClick={this.handlePaging}>Next Week</button>
		</div>;
	}
}

const mapStateToProps = state => {
	return {
		days: state.Days,
		anchorDate: state.Days.length ? dateFromSlot(state.Days[0][0]) : new Date()
	};
};


const mapDispatchToProps = dispatch => {
	return {
		fetchWeek: (anchorDate, direction) => dispatch(fetchWeek(store.getState(), anchorDate, direction))
	};
};

const Schedule = connect(mapStateToProps, mapDispatchToProps)(ConnectedSchedule);

export default Schedule;