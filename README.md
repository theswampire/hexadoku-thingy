# README for a hastily hacked hexadoku hinter

Well, it's a tool that I hacked in one day to give hints for what values are possible in a hexadoku. 
Don't try creating sudokus with sizes other than 9 and 16, it'll probably just give you something weird.
And the values are always starting from 0. So if you try it for normal sudokus, you'll have to add 1 to all values :)

What I learned or realized: 
- I think I just became a baby gopher! Why didn't I try Go earlier
- I realized that using react for such things is kind of / sort of / definitely a mentally underperforming decision. Why would anyone want to manage state in 2 different places?

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
