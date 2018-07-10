import React, { Component } from 'react';
import { connect } from 'react-redux';
import Day from './Day.jsx';
import { dateFromSlot } from './store/funcs';
import { fetchWeek } from './store/schedule';

const Schedule = ({ anchorDate, days, fetchWeek }) => {
	const onPaginate = e => fetchWeek(anchorDate, e.target.className);

	return (
		<div>
			<button className="prev" onClick={onPaginate}>Previous Week</button>
			<div id="week">
				{days.map((day, index) => <Day key={index} slots={day} />)}
			</div>
			<button className="next" onClick={onPaginate}>Next Week</button>
		</div>
	);
};

const mapStateToProps = state => {
	const days = state.Days;

	return {
		days,
		anchorDate: days.length ? dateFromSlot(days[0][0]) : new Date(),
	};
};

const mapDispatchToProps = { fetchWeek };

const ConnectedSchedule = connect(mapStateToProps, mapDispatchToProps)(Schedule);

export default ConnectedSchedule;
