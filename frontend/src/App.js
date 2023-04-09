// import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import Releases from './pages/releases';
import Projects from './pages/projects';
import Github from './pages/github';

function App() {

  return (
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path='/release-notes' element={<Releases/>} />
          <Route path='/projects' element={<Projects/>} />
          <Route path='/github' element={<Github/>} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;