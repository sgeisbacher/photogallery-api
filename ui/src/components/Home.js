import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';


class Home extends Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    fetch('http://localhost:8080/labels')
      .then((resp) => resp.json())
      .then((labels) => this.setState({labels}))
  }

  render() {
    if (this.state.labels) {
      return (
        <div>
          <ul>
            {this.state.labels.map((label) =>
              <li key={label.ID}>
                <NavLink to={`/label/${label.ID}`}>{label.Name}</NavLink>
              </li>)}
          </ul>
        </div>
      )
    }
    return <div>loading...</div>
  }
}

export default Home