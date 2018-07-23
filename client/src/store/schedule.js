import fetch from 'cross-fetch';

const SERVER = process.env.REACT_APP_SERVER_ENDPOINT;
const WEEK_ENDPOINT = '/schedule/week.json';

const requestWeek = () => ({ type: 'REQUEST_WEEK' });
const receiveWeek = (anchorDate, json) => ({ type: 'RECEIVE_WEEK', anchorDate, data: json });

export const fetchWeek = (anchorDate, direction = '') => {
	return (dispatch, getState) => {
		const { Company, ClientID } = getState();

		dispatch(requestWeek());

		const url = (
			SERVER +
			WEEK_ENDPOINT +
			'?companyID=' + Company.ID +
			'&clientID=' + ClientID +
			'&year=' + anchorDate.getFullYear() +
			'&month=' + (anchorDate.getMonth() + 1) + // javascript is zero-based, go is one-based
			'&day=' + anchorDate.getDate() +
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
