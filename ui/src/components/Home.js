import React, { Component } from 'react';
import { NavLink } from 'react-router-dom';

class Home extends Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    fetch('/labels')
      .then((resp) => resp.json())
      .then((labels) => this.setState({labels}))
  }

  render() {
    const alphabeticalSorter = (a, b) => a.Name > b.Name
    if (this.state.labels) {
      return (
        <div>
          <ul>
            {this.state.labels.slice().sort(alphabeticalSorter).map((label) =>
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