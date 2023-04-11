// import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';
import ModalWarn from './components/ModalWarn';

import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import ReleasesNotes from './pages/release-notes';
import Projects from './pages/projects';
// import TestingRoute from './pages/testing-route';
import Github from './pages/github';

function App() {

  return (
    <>
    <Router>
      <div>
        <NavBar />
        <Routes>
          {/* <Route path='/release-notes' element={<ReleasesNotes/>} /> */}
          <Route path='/projects' element={<Projects/>} />
          <Route path='/github' element={<Github/>} />
          <Route path='/notifications' element={<Github/>} />
          <Route path="/repo/:owner/:repo/releases" element={<ReleasesNotes/>} />
        </Routes>
      </div>
    </Router>

    {/* show message when redis is down */}
    {/* <ModalWarn /> */}

    </>

  );
}

export default App;