// import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';
// import ModalWarn from './components/ModalWarn';

import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import ReleasesNotes from './pages/release-notes';
import Projects from './pages/projects';
import PotentialUpdates from './pages/potential-updates';
import WatcherList from './pages/watcherList';
import APIStats from './pages/apiStats';

function App() {

  return (
    <>
    <Router>
      <div>
        <NavBar />
        <Routes>
          {/* <Route path='/release-notes' element={<ReleasesNotes/>} /> */}
          <Route path='/projects' element={<Projects/>} />
          <Route path='/watcher-list' element={<WatcherList/>} />
          <Route path='/notifications' />
          <Route path="/repo/:owner/:repo/releases" element={<ReleasesNotes/>} />
          <Route path="/api-stats" element={<APIStats/>} />
          <Route path="/potential-updates" element={<PotentialUpdates/>} />
        </Routes>
      </div>
    </Router>

    {/* show message when redis is down */}
    {/* <ModalWarn /> */}

    </>

  );
}

export default App;