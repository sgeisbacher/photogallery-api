import React, { Component } from 'react';

class Media extends Component {

  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
  }

  render() {
    const id = this.props.match.params.id
    return (
      <div>
        <div>
          <div className="mediaContainer">
            <div style={{ display: 'inline-block', margin: '5px' }}><img src={`http://localhost:8080/data/media/big/${id}`}/></div>
          </div>
        </div>
      </div>
    )
  }
}

export default Media