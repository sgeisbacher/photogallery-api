import React, { Component } from 'react';

class Label extends Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    fetch(`http://localhost:8080/labels/${this.props.match.params.id}`)
      .then((resp) => resp.json())
      .then((label) => this.setState({label}))
  }

  render() {
    return (
      <div>
        { this.state.label ? 
          <div>
            <h1>Label '{this.state.label.Name}'</h1> 
            <hr/>
            <div className="thumbContainer">
              {this.state.label.Medias.map((media) => (
                <div style={{ display: 'inline-block', margin: '5px' }}><img src={`http://localhost:8080${media.ThumbUrl}`}/></div>
              ))}
            </div>
          </div>
          : <span>loading ...</span>
          }
      </div>
    )
  }
}

export default Label