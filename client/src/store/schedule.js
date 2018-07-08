import fetch from 'cross-fetch';

const SERVER = 'http://35.192.118.156:8080';
const WEEK_ENDPOINT = '/schedule/week.json';

function requestWeek() {
	return {
		type: 'REQUEST_WEEK'
	};
}

function receiveWeek(anchorDate, json) {
	return {
		type: 'RECEIVE_WEEK',
		data: json
	};
}

export const fetchWeek = (state, anchorDate, direction='') => {
	return (dispatch) => {
		dispatch(requestWeek());

		const companyID = state.Company.ID;
		const clientID = state.ClientID;
		const year = anchorDate.getFullYear();
		const month = anchorDate.getMonth() + 1; // javascript is zero-based, go is one-based
		const day = anchorDate.getDate();
		const url = (
			SERVER +
			WEEK_ENDPOINT +
			'?companyID=' + companyID +
			'&clientID=' + clientID +
			'&year=' + year +
			'&month=' + month +
			'&day=' + day +
			'&direction=' + direction
		);
		return fetch(url).then(
			res => res.json(),
			// Do not use catch, because that will also catch
			// any errors in the dispatch and resulting render,
			// causing a loop of 'Unexpected batch number' errors.
			// https://github.com/facebook/react/issues/6895
			err => console.error('error in fetchWeek', err)
		).then(json => dispatch(receiveWeek(anchorDate, json)));
	};
};