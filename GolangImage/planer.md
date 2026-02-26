## header slime

- Ok basiccly this will be and cli app the will take frames frm opencvs python camera send the bits the a golang 
- program via tcp and then the golang app will print the image frames in ascii and will update every couple of seconds in order to make 
 it look like the camera is being dispayed in the terminal there will be some challanges with speading up conversion but will face that when we 
	get there for now my goal is
1. Get python to send bytes of a frame in image to a golang tcp server and just display the bytes
2. after recinving said bytes then use the previos image rendering program in order to display the single picture in the terminal

start with this for now
