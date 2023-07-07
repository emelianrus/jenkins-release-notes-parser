
import 'bootstrap/dist/css/bootstrap.min.css';

import React, { useState, useEffect } from "react";
import PluginChangesCard from '../components/PluginChangesCard';
import PotentialUpdatesList from '../components/PotentialUpdatesList';

// TODO: add loading status during long check-deps call
function PluginChanges() {

  const [pluginsDiff, setPluginsDiff] = useState([]);
  const [projects, setProjects] = useState([]);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {

    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/check-deps`);
      const data = await response.json();
      console.log(data)
    } catch (error) {
      console.error(error);
    }

    try {
      const response = await fetch(`http://localhost:8080/plugin-manager/get-fixed-deps-diff`);
      const data = await response.json();
      setPluginsDiff(data);

      setProjects(data);
      // console.log(data);
    } catch (error) {
      console.error(error);
    }

  //   setProjects([
  //     {
  //       "Name": "kubernetes-client-api",
  //       "CurrentVersion": "5.4.2",
  //       "NewVersion": "5.12.1-187.v577c3e368fb_6",
  //       "ReleaseNotes": [
  //               {
  //                   "Name": "6.4.1-215.v2ed17097a_8e9",
  //                   "Tag": "6.4.1-215.v2ed17097a_8e9",
  //                   "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse okhttp-api plugin instead of bundling okhttp from kubernetes-client (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/179\"\u003e#179\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse okhttp-api plugin instead of bundling okhttp from kubernetes-client (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/179\"\u003e#179\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump git-changelist-maven-extension from 1.4 to 1.6 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/178\"\u003e#178\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eRoutine updates (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/174\"\u003e#174\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                   "HTMLURL": "",
  //                   "CreatedAt": "2023-03-02T12:43:42Z"
  //               },
  //               {
  //                   "Name": "6.4.1-208.vfe09a_9362c2c",
  //                   "Tag": "6.4.1-208.vfe09a_9362c2c",
  //                   "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 6.3.1 to 6.4.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/175\"\u003e#175\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 6.3.1 to 6.4.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/175\"\u003e#175\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                   "HTMLURL": "",
  //                   "CreatedAt": "2023-02-16T13:26:16Z"
  //               },
  //               {
  //                   "Name": "6.3.1-206.v76d3b_6b_14db_b",
  //                   "Tag": "6.3.1-206.v76d3b_6b_14db_b",
  //                   "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"breaking-changes\"\u003e❗ Breaking changes\u003c/h2\u003e\n\n\u003cp\u003e\u003cstrong\u003eVersion 6.x of kubernetes-client has breaking changes, please hold on upgrading to this release until other plugins in the Jenkins ecosystem have been released with compatibility changes.\u003c/strong\u003e\nCompatible versions\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-credentials-plugin/releases/tag/kubernetes-credentials-0.10.0\" target=\"_blank\"\u003ekubernetes-credentials 0.10.0\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-cli-plugin/releases/tag/kubernetes-cli-1.11.0\" target=\"_blank\"\u003ekubernetes-cli 1.11.0\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-plugin/releases/tag/3802.vb_b_600831fcb_3\" target=\"_blank\"\u003ekubernetes 3802.vb_b_600831fcb_3\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-credentials-provider-plugin/releases/tag/1.208.v128ee9800c04\" target=\"_blank\"\u003ekubernetes-credentials-provider 1.208.v128ee9800c04\u003c/a\u003e\u003c/p\u003e\n\n\u003ch2 id=\"other-changes\"\u003e✍ Other changes\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003echore: use jenkins infra maven cd reusable workflow (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/150\"\u003e#150\u003c/a\u003e) \u003ca href=\"https://github.com/jetersen\"\u003e@jetersen\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.2 to 6.3.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/163\"\u003e#163\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump git-changelist-maven-extension from 1.3 to 1.4 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/155\"\u003e#155\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump plugin from 4.40 to 4.50 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/165\"\u003e#165\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                   "HTMLURL": "",
  //                   "CreatedAt": "2023-01-09T09:56:31Z"
  //               },
  //               {
  //                   "Name": "5.12.2-193.v26a_6078f65a_9",
  //                   "Tag": "5.12.2-193.v26a_6078f65a_9",
  //                   "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.1 to 5.12.2 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/149\"\u003e#149\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"maintenance\"\u003e👻 Maintenance\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse Eclipse Temurin, not AdoptOpenJDK in action (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/146\"\u003e#146\u003c/a\u003e) \u003ca href=\"https://github.com/MarkEWaite\"\u003e@MarkEWaite\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.1 to 5.12.2 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/149\"\u003e#149\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump plugin from 4.37 to 4.40 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/148\"\u003e#148\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                   "HTMLURL": "",
  //                   "CreatedAt": "2022-05-18T15:35:55Z"
  //               }
  //           ],
  //       "Type": 2
  //   },
  //   {
  //     "Name": "kubernetes-client-api222222",
  //     "CurrentVersion": "5.4.2",
  //     "NewVersion": "5.12.1-187.v577c3e368fb_6",
  //     "ReleaseNotes":
  //         [
  //             {
  //                 "Name": "6.4.1-215.v2ed17097a_8e9",
  //                 "Tag": "6.4.1-215.v2ed17097a_8e9",
  //                 "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse okhttp-api plugin instead of bundling okhttp from kubernetes-client (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/179\"\u003e#179\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse okhttp-api plugin instead of bundling okhttp from kubernetes-client (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/179\"\u003e#179\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump git-changelist-maven-extension from 1.4 to 1.6 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/178\"\u003e#178\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eRoutine updates (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/174\"\u003e#174\u003c/a\u003e) \u003ca href=\"https://github.com/Vlatombe\"\u003e@Vlatombe\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                 "HTMLURL": "",
  //                 "CreatedAt": "2023-03-02T12:43:42Z"
  //             },
  //             {
  //                 "Name": "6.4.1-208.vfe09a_9362c2c",
  //                 "Tag": "6.4.1-208.vfe09a_9362c2c",
  //                 "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 6.3.1 to 6.4.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/175\"\u003e#175\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 6.3.1 to 6.4.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/175\"\u003e#175\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                 "HTMLURL": "",
  //                 "CreatedAt": "2023-02-16T13:26:16Z"
  //             },
  //             {
  //                 "Name": "6.3.1-206.v76d3b_6b_14db_b",
  //                 "Tag": "6.3.1-206.v76d3b_6b_14db_b",
  //                 "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"breaking-changes\"\u003e❗ Breaking changes\u003c/h2\u003e\n\n\u003cp\u003e\u003cstrong\u003eVersion 6.x of kubernetes-client has breaking changes, please hold on upgrading to this release until other plugins in the Jenkins ecosystem have been released with compatibility changes.\u003c/strong\u003e\nCompatible versions\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-credentials-plugin/releases/tag/kubernetes-credentials-0.10.0\" target=\"_blank\"\u003ekubernetes-credentials 0.10.0\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-cli-plugin/releases/tag/kubernetes-cli-1.11.0\" target=\"_blank\"\u003ekubernetes-cli 1.11.0\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-plugin/releases/tag/3802.vb_b_600831fcb_3\" target=\"_blank\"\u003ekubernetes 3802.vb_b_600831fcb_3\u003c/a\u003e\n* ✅ \u003ca href=\"https://github.com/jenkinsci/kubernetes-credentials-provider-plugin/releases/tag/1.208.v128ee9800c04\" target=\"_blank\"\u003ekubernetes-credentials-provider 1.208.v128ee9800c04\u003c/a\u003e\u003c/p\u003e\n\n\u003ch2 id=\"other-changes\"\u003e✍ Other changes\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003echore: use jenkins infra maven cd reusable workflow (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/150\"\u003e#150\u003c/a\u003e) \u003ca href=\"https://github.com/jetersen\"\u003e@jetersen\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.2 to 6.3.1 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/163\"\u003e#163\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump git-changelist-maven-extension from 1.3 to 1.4 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/155\"\u003e#155\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump plugin from 4.40 to 4.50 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/165\"\u003e#165\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                 "HTMLURL": "",
  //                 "CreatedAt": "2023-01-09T09:56:31Z"
  //             },
  //             {
  //                 "Name": "5.12.2-193.v26a_6078f65a_9",
  //                 "Tag": "5.12.2-193.v26a_6078f65a_9",
  //                 "BodyHTML": "\u003c!-- Optional: add a release summary here --\u003e\n\n\u003ch2 id=\"new-features-and-improvements\"\u003e🚀 New features and improvements\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.1 to 5.12.2 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/149\"\u003e#149\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"maintenance\"\u003e👻 Maintenance\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eUse Eclipse Temurin, not AdoptOpenJDK in action (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/146\"\u003e#146\u003c/a\u003e) \u003ca href=\"https://github.com/MarkEWaite\"\u003e@MarkEWaite\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n\n\u003ch2 id=\"dependency-updates\"\u003e📦 Dependency updates\u003c/h2\u003e\n\n\u003cul\u003e\n\u003cli\u003eBump revision from 5.12.1 to 5.12.2 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/149\"\u003e#149\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003cli\u003eBump plugin from 4.37 to 4.40 (\u003ca href=\"https://github.com/jenkinsci/plugin-installation-manager-tool/pull/148\"\u003e#148\u003c/a\u003e) \u003ca href=\"https://github.com/dependabot\"\u003e@dependabot\u003c/a\u003e\u003c/li\u003e\n\u003c/ul\u003e\n",
  //                 "HTMLURL": "",
  //                 "CreatedAt": "2022-05-18T15:35:55Z"
  //             }
  //         ],
  //     "Type": 2
  // }
  //   ])


    // try {
    //   // TODO: https://api.github.com/repos/OWNER/REPO/releases
    //   // repos/:owner/:repo/releases
    //   // /project/:owner/:repo/releases
    //   const response = await fetch(`http://localhost:8080/potential-updates`);
    //   const data = await response.json();

    //   // pass as new single object instead of several params
    //   setProjects(data);
    // } catch (error) {
    //   console.error(error);
    // }

  };

  const pluginsArray = [];
  for (const key in pluginsDiff) {
    pluginsArray.push({
      key,
      project: pluginsDiff[key]
    });
  }

  const pluginCards = pluginsArray.map(plugin => (
    <PluginChangesCard key={plugin.key} project={plugin.project} />
  ));
  return (
    <div>
      <div className="project-list">
      <div className="container-sm mt-5 ml-5 mr-5">
        <h3>Title</h3>
        <div className="table-responsive">
          <table className="table">
            <thead className="thead-light">
              <tr>
                <th>Project</th>
                <th>From version</th>
                <th>To version</th>
                <th>Type</th>
              </tr>
            </thead>

            {/* {projectCards.map(project => <PluginManagerCard key={project.key} project={project} />)} */}
            {pluginsDiff === undefined
              ? <tbody><tr><td>No projects to display</td></tr></tbody>
              : pluginCards
            }

          </table>
          <b>RELEASE NOTES</b>
          <PotentialUpdatesList projects={projects}/>
        </div>
      </div>
    </div>
    </div>
  );
}

export default PluginChanges;