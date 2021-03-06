import React, { Component } from 'react';
import { Link } from 'react-router-dom';

class Label extends Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    fetch(`/labels/${this.props.match.params.id}`)
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
                <div style={{ display: 'inline-block', margin: '5px' }}>
                  <Link to={`/media/${media.Hash}`}><img src={media.ThumbUrl}/></Link> 
                </div>
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