package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	color = mgl32.Vec3{1, 0, 0} // vermelho
)

const (
	width  = 800
	height = 600
	title  = "LearnOpenGL"
)

// A função init() é executada antes da main().
// É uma convenção em Go para garantir que a thread principal seja a do runtime.
func init() {
	// A GLFW exige que ela seja executada na thread principal do sistema.
	runtime.LockOSThread()
}

func main() {
	// Inicializa a GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("Falha ao inicializar o GLFW:", err)
	}

	// Garante que a GLFW seja encerrada no final da função.
	defer glfw.Terminate()

	// Solicita um contexto OpenGL 4.6.
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	// Solicita o perfil de compatibilidade, que inclui as funções "legacy".
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)
	// Socilita a remoção das funções legacy
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.False)

	// Cria a janela
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatalln("Falha ao criar a janela:", err)
	}

	// Define a janela atual
	window.MakeContextCurrent()

	// Registra a função de callback para o redimensionamento
	window.SetSizeCallback(resizeCallback)

	// Registra a função de callback do teclado.
	window.SetKeyCallback(keyboardCallback)

	// Inicializa o OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalln("Falha ao iniciar o OpenGL:", err)
	}

	// Exibe a versão do OpenGL que está sendo usada.
	log.Printf("Versão do OpenGL: %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// Loop principal da janela enquanto a janlea ta aberta
	for !window.ShouldClose() {
		// Limpa a tela com uma cor de fundo preto
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.Color3f(color.X(), color.Y(), color.Z()) // Define a cor vermelha
		gl.Begin(gl.TRIANGLES)                      // Inicia um triângulo
		gl.Vertex2f(-0.5, -0.5)                     // canto inferior esquerdo
		gl.Vertex2f(0.5, -0.5)                      // canto inferior direito
		gl.Vertex2f(0, 0.5)                         // meio em cima
		gl.End()                                    // Termina o desenho

		// Troca os buffers de vídeo. É essencial para a renderização de gráficos.
		window.SwapBuffers()

		// Processa todos os eventos somente quando disparados
		glfw.WaitEvents()
	}
}

// resizeCallback é chamada toda vez que a janela é redimensionada.
func resizeCallback(w *glfw.Window, width int, height int) {
	// gl.Viewport define a área de renderização do OpenGL.
	// O primeiro e segundo parâmetros são as coordenadas (x, y) do canto inferior esquerdo.
	// O terceiro e quarto parâmetros são a largura e altura da nova área.
	gl.Viewport(0, 0, int32(width), int32(height))
}

// keyboardCallback é chamada sempre que uma tecla é pressionada, liberada ou repetida.
func keyboardCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// Se a tecla ESC (escape) for pressionada, a janela deve ser fechada.
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}

	if key == glfw.KeyR && action == glfw.Press {
		// Define as cores para comparação.
		red := mgl32.Vec3{1, 0, 0}
		green := mgl32.Vec3{0, 1, 0}
		blu := mgl32.Vec3{0, 0, 1}

		if color.ApproxEqual(red) {
			color = blu
		} else if color.ApproxEqual(blu) {
			color = green
		} else {
			color = red
		}
	}
}
