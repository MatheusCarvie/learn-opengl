package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sqweek/dialog"
)

var color = mgl32.Vec3{1, 0, 0} // vermelho
var cubPosition = mgl32.Vec2{0, 0}
var cubeSpeed float32 = 150

var (
	lastFrameTime float64
	deltaTime     float64
	fps           float64
	frames        int
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
		dialog.Message("Falha ao inicializar o GLFW: %s", err).Title("Erro").Error()
		return
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
		dialog.Message("Falha ao criar a janela: %s", err).Title("Erro").Error()
		return
	}

	// Define a janela atual
	window.MakeContextCurrent()

	// Desabilita o V-Sync fazendo o fps nao ser fixo aos hz do monitor
	glfw.SwapInterval(0)

	// Registra a função de callback para o redimensionamento
	window.SetSizeCallback(resizeCallback)

	// Registra a função de callback do teclado.
	window.SetKeyCallback(keyboardCallback)

	// Inicializa o OpenGL
	if err := gl.Init(); err != nil {
		dialog.Message("Falha ao iniciar o OpenGL: %s", err).Title("Erro").Error()
		return
	}

	// Exibe a versão do OpenGL que está sendo usada.
	// dialog.Message("Versão do OpenGL: %s", gl.GoStr(gl.GetString(gl.VERSION))).Title("Informação").Info()

	// Inicializa o tempo do último frame.
	lastFrameTime = glfw.GetTime()

	// Loop principal da janela enquanto a janlea ta aberta
	for !window.ShouldClose() {
		// Calcula o FPS
		getFramePerSeconds()

		// Limpa a tela com uma cor de fundo preto
		gl.ClearColor(0, 0, 0, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.PointSize(40)

		moveCube(window)
		gl.LoadIdentity() // reseta posição e rotação
		gl.Translatef(cubPosition.X(), cubPosition.Y(), 0)

		gl.Color3f(color.X(), color.Y(), color.Z()) // Define a cor vermelha
		gl.Begin(gl.POINTS)
		gl.Vertex2f(0.0, 0.0)
		gl.End() // Termina o desenho

		// Troca os buffers de vídeo. É essencial para a renderização de gráficos.
		window.SwapBuffers()

		// Processa todos os eventos somente quando disparados
		glfw.PollEvents()
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

func getFramePerSeconds() {
	// Obtém o tempo do frame atual.
	currentTime := glfw.GetTime()

	// Calcula o Delta Time (o tempo que o frame anterior levou para renderizar)
	deltaTime = currentTime - lastFrameTime

	// Atualiza o tempo do último frame.
	lastFrameTime = currentTime

	// Calcula o FPS
	frames++
	if deltaTime > 0.0 {
		// A cada segundo, atualiza o FPS.
		fps = float64(frames) / deltaTime
		// Loga o FPS a cada segundo, ou exibe na janela.
		fmt.Printf("FPS: %.2f\n", fps)
		frames = 0
	}
}

func moveCube(window *glfw.Window) {
	// dentro do loop principal, antes de desenhar
	var cubeSpeed float32 = 1.0 // unidades por segundo

	// Calcula direção pura (sem aplicar velocidade ainda)
	direction := mgl32.Vec2{0, 0}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		direction = direction.Add(mgl32.Vec2{0, 1})
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		direction = direction.Sub(mgl32.Vec2{0, 1})
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		direction = direction.Add(mgl32.Vec2{1, 0})
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		direction = direction.Sub(mgl32.Vec2{1, 0})
	}

	// Se tiver input, normaliza
	if direction.Len() > 0 {
		direction = direction.Normalize()
	}

	cubPosition = cubPosition.Add(direction.Mul(float32(deltaTime) * cubeSpeed))
}
