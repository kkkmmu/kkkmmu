Components
The model is the central component of the pattern. It is the application's dynamic data structure, independent of the user interface.[6] It directly manages the data, logic and rules of the application.
A view can be any output representation of information, such as a chart or a diagram. Multiple views of the same information are possible, such as a bar chart for management and a tabular view for accountants.
The third part or section, the controller, accepts input and converts it to commands for the model or view.[7]

Interactions
In addition to dividing the application into three kinds of components, the model–view–controller design defines the interactions between them.[8]

The model is responsible for managing the data of the application. It receives user input from the controller.
The view means presentation of the model in a particular format.
The controller responds to the user input and performs interactions on the data model objects. The controller receives the input, optionally validates it and then passes the input to the model.
