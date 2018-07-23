import { createStore, applyMiddleware } from "redux";
import thunkMiddleware from 'redux-thunk'
import rootReducer from "./reducers";
import { fetchWeek } from "./schedule";

const store = createStore(
	rootReducer,
	applyMiddleware(
		thunkMiddleware
	)
);

store.dispatch(fetchWeek(new Date()));

export default store;
