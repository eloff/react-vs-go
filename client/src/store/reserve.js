import fetch from 'cross-fetch';
import store from './index';
import { dateFromSlot } from './funcs';

const SERVER = process.env.REACT_APP_SERVER_ENDPOINT;
const RESERVATION_ENDPOINT = '/schedule/reserve';

const requestReservation = () => ({ type: 'REQUEST_RESERVATION' });
const madeReservation = data => ({ type: 'MADE_RESERVATION', data });

const completeReservation = (state, date) => ({
	...state,
	Days: state.Days.map(week => {
		if (dateFromSlot(week[0]).getDate() !== date.getDate()) {
			return [...week];
		}

		return week.map(day => {
			if (dateFromSlot(day).getTime() !== date.getTime()) {
				return { ...day };
			}

			return {
				...day,
				ClientID: state.ClientID,
				Available: false,
			};
		});
	}),
});

export const makeReservation = date => {
	return (dispatch, getState) => {
		const { Company, ClientID } = getState();

		dispatch(requestReservation());

		const url = (
			SERVER +
			RESERVATION_ENDPOINT +
			'?companyID=' + Company.ID +
			'&clientID=' + ClientID +
			'&year=' + date.getFullYear() +
			'&month=' + (date.getMonth() + 1) + // javascript is zero-based, go is one-based
			'&day=' + date.getDate() +
			'&hour=' + date.getHours() +
			'&minute=' + date.getMinutes()
		);

		return fetch(url).then(
			res => dispatch(madeReservation(completeReservation(getState(), date))),
			err => console.error('error in makeReservation', err)
		);
	};
};
