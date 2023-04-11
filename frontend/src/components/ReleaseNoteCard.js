

function ReleaseNoteCard({ project }) {
    return (
      // TODO: fix class name
      <div className="cardd">
        <section aria-labelledby="hd-0de31afa">
            <div className="d-flex flex-column flex-md-row my-5 flex-justify-center">
              {/* left menu */}
              <div className="col-md-2 d-flex flex-md-column flex-row flex-wrap pr-md-6 mb-2 mb-md-0 flex-items-start pt-md-4">
                <div className="mb-2 f4 mr-3 mr-md-0 col-12">
                  {project.CreatedAt}
                </div>
                <div className="mr-3 mr-md-0 d-flex" data-pjax="#repo-content-pjax-container" data-turbo-frame="repo-content-turbo-frame">
                  {project.Tag}
                </div>
              </div>
              {/* left menu end */}

              <div className="col-md-9">
                <div data-view-component="true" className="Box">
                    <div data-view-component="true" className="Box-body">
                      <div className="d-flex flex-md-row flex-column">
                        <div className="d-flex flex-row flex-1 mb-3 wb-break-word">
                          <div className="flex-1" data-pjax="#repo-content-pjax-container" data-turbo-frame="repo-content-turbo-frame">
                            <span data-view-component="true" className="f1 text-bold d-inline mr-3">
                              <a href="/#" data-view-component="true" className="custom-link" data-turbo-frame="repo-content-turbo-frame">{project.Name}</a>
                            </span>
                          </div>
                        </div>
                      </div>

                      <div data-pjax="true" data-test-selector="body-content" data-view-component="true" className="markdown-body my-3">
                        {/* here markdown */}
                        <div dangerouslySetInnerHTML={{ __html: project.BodyHTML }} />
                        {/* here markdown end*/}
                      </div>
                    </div>
                  </div>
              </div>
            </div>
          </section>
      </div>
    );
}

export default ReleaseNoteCard;


