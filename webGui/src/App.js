import ReactDOM from 'react-dom';
import React from 'react';


function Welcome(props){
  return <h1>hello world,{props.name}</h1>
}

function App() {
  return (
    <div>
      <Welcome name="Sara" />
      <Welcome name="Cahal" />
      <Welcome name="Edite" />
    </div>
  );
}

ReactDOM.render(
  <App />,
  document.getElementById('root')
);
export default App;