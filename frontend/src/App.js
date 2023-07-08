
import 'bootstrap/dist/css/bootstrap.min.css';

import NavBar from './components/NavBar';

import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';

import PluginManager from './pages/plugin-manager';
import PluginChanges from './pages/plugin-changes';
import WatcherList from './pages/watcherList';

function App() {

  return (
    <>
    <Router>
      <div>
        <NavBar />
        <Routes>
          <Route path='/watcher-list' element={<WatcherList/>} />
          <Route path="/plugin-manager" element={<PluginManager/>} />
          <Route path="/plugin-changes" element={<PluginChanges/>} />
        </Routes>
      </div>
    </Router>
    </>

  );
}

export default App;