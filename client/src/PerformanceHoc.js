import React from 'react';

export class PerformantCall extends React.Component {
    startTimer = measurementName => {
        this.start = performance.now();
        this.measurementName = measurementName;
    }

    componentDidUpdate(prevProps) {
        this.end = performance.now();
        console.info(`${this.measurementName} took ${this.end - this.start}ms`);
    }

    render() {
        const { Component, ...props } = this.props;

        return <Component {...props} onStartPerformanceTimer={this.startTimer} />
    }
};

const PerformanceHoc = Component => {
    return props => <PerformantCall Component={Component} {...props} />;
};

export default PerformanceHoc;
