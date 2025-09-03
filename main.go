package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"

	"github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sqweek/dialog"
)

var snackPosition = mgl32.Vec2{0, 0}
var enemyPosition = mgl32.Vec2{0, 0}
var tailsPosition = []mgl32.Vec2{}
var points = 0

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
	glfw.SwapInterval(1)

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

	// Limpa a tela com uma cor de fundo preto
	gl.ClearColor(0, 0, 0, 1)
	gl.PointSize(40)

	spawnEnemy()

	// Loop principal da janela enquanto a janlea ta aberta
	for !window.ShouldClose() {
		// Calcula o FPS
		getFramePerSeconds()

		gl.Clear(gl.COLOR_BUFFER_BIT) // Limpa os buffer
		gl.LoadIdentity()             // reseta posição e rotação

		moveSnack(window) // Movimenta o player

		renderEnemy() // Renderiza o inimigo
		renderSnack() // Renderiza o player
		renderTail()  // Renderiza a cauda

		// Checka colisao entre o player e o inimigo
		if checkCollision(snackPosition, enemyPosition, 0.15) {
			spawnEnemy()
			points++
			fmt.Println("Pontos:", points)
		}

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
		// fmt.Printf("FPS: %.2f\n", fps)
		frames = 0
	}
}

func moveSnack(window *glfw.Window) {
	// dentro do loop principal, antes de desenhar
	var snackSpeed float32 = 1.0 // unidades por segundo

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

	// Adiciona a cauda com a posição do player no index 0
	tailsPosition = append([]mgl32.Vec2{snackPosition}, tailsPosition...)
	if points > 0 && len(tailsPosition) > points {
		// Remove caudas que passar dos pontos
		tailsPosition = tailsPosition[:points]
	}

	// Salva a nova posição
	newPosition := snackPosition.Add(direction.Mul(float32(deltaTime) * snackSpeed))

	// Limita a posição para nao sair da tela
	snackPosition = mgl32.Vec2{
		clamp(newPosition.X(), -1, 1),
		clamp(newPosition.Y(), -1, 1),
	}
}

func renderSnack() {
	gl.Color3f(0, 1, 0) // Define a cor verde
	gl.Begin(gl.POINTS)
	gl.Vertex2f(snackPosition.X(), snackPosition.Y())
	gl.End() // Termina o desenho
}

func renderTail() {
	for _, pos := range tailsPosition {
		gl.Color3f(0, 1, 0)
		gl.Begin(gl.POINTS)
		gl.Vertex2f(pos.X(), pos.Y())
		gl.End()
	}
}

func renderEnemy() {
	gl.Color3f(1, 0, 0) // Define a cor vermelha
	gl.Begin(gl.POINTS)
	gl.Vertex2f(enemyPosition.X(), enemyPosition.Y())
	gl.End() // Termina o desenho
}

func spawnEnemy() {
	randomX := randomFloat32(-1, 1)
	randomY := randomFloat32(-1, 1)
	enemyPosition = mgl32.Vec2{randomX, randomY}
}

// Função para gerar um número aleatório entre min e max
func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// Colisão ocorre quando a distância entre os centros é menor que a soma dos raios
func checkCollision(pos1 mgl32.Vec2, pos2 mgl32.Vec2, radius float32) bool {
	dx := float64(pos1.X() - pos2.X())
	dy := float64(pos1.Y() - pos2.Y())
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < float64(radius)
}

func clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
